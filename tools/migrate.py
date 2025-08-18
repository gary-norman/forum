import sqlite3

src_path = "old.db"
tgt_path = "dev_forum_database.db"

# Connect to source DB
src = sqlite3.connect(src_path)
src.execute(f"ATTACH DATABASE '{tgt_path}' AS target")
src.execute("PRAGMA foreign_keys = OFF")

# Get table list in creation order (so parent tables before children)
tables = [
    row[0]
    for row in src.execute(
        "SELECT name FROM sqlite_master WHERE type='table' AND name NOT LIKE 'sqlite_%' ORDER BY name"
    )
]

# Store ID remapping for each table where we skip PK
id_map = {}

for table in tables:
    # Find PK column
    pk_col = None
    cols_info = list(src.execute(f"PRAGMA table_info({table})"))
    for col in cols_info:
        if col[5] == 1:  # pk flag
            pk_col = col[1]
            break

    if table.lower() == "users":
        # Preserve IDs exactly
        cols = [col[1] for col in cols_info]
        src.execute(
            f"INSERT OR IGNORE INTO target.{table} ({', '.join(cols)}) "
            f"SELECT {', '.join('main.' + table + '.' + c for c in cols)} "
            f"FROM main.{table}"
        )
    else:
        # Skip PK and map old -> new IDs
        non_pk_cols = [col[1] for col in cols_info if col[1] != pk_col]
        placeholders = ", ".join("?" for _ in non_pk_cols)
        select_sql = f"SELECT {', '.join(non_pk_cols)}, {pk_col} FROM main.{table}"

        for row in src.execute(select_sql):
            *non_pk_values, old_id = row
            cur = src.execute(
                f"INSERT OR IGNORE INTO target.{table} ({', '.join(non_pk_cols)}) "
                f"VALUES ({placeholders})",
                non_pk_values,
            )
            new_id = cur.lastrowid
            id_map.setdefault(table, {})[old_id] = new_id

# Now fix foreign keys in child tables
for table in tables:
    fk_list = list(src.execute(f"PRAGMA foreign_key_list({table})"))
    if not fk_list:
        continue

    # Check if this table contains FK columns that need updating
    for fk in fk_list:
        fk_col = fk[3]  # child column name
        parent_table = fk[2]
        if parent_table in id_map:
            # Update rows to reflect new PK values from parent table
            for old_id, new_id in id_map[parent_table].items():
                src.execute(
                    f"UPDATE target.{table} SET {fk_col} = ? WHERE {fk_col} = ?",
                    (new_id, old_id),
                )

src.commit()
src.execute("PRAGMA foreign_keys = ON")
print("Migration complete with foreign keys preserved.")
