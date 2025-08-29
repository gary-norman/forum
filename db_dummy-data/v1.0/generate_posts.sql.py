import random

users = [
    (
        "dreamsofcode",
        "5f90bcae-b649-4b26-8f6c-1fa9217303f6",
        "doc@codex.com",
        "doc.jpg",
        "$2a$14$Azk3S9Zkbt0LrV/uN8WVMuoFjl0MJxpuXp31ps6lyChyQTL2xALD6",
    ),
    (
        "webDevSimplified",
        "94e52033-8ec2-413d-aa49-7a69d23f6836",
        "wds@codex.com",
        "wds.png",
        "$2a$14$Azk3S9Zkbt0LrV/uN8WVMuoFjl0MJxpuXp31ps6lyChyQTL2xALD6",
    ),
    (
        "frontendfriends",
        "3da1bc7e-e9bd-4518-8151-6d5ff18d06c3",
        "fef@codex.com",
        "kpowell.jpeg",
        "$2a$14$Azk3S9Zkbt0LrV/uN8WVMuoFjl0MJxpuXp31ps6lyChyQTL2xALD6",
    ),
    (
        "kamil",
        "b920bfdc-7a44-445d-ba09-a665df542ac9",
        "kamil@codex.com",
        "kamil.png",
        "$2a$14$Azk3S9Zkbt0LrV/uN8WVMuoFjl0MJxpuXp31ps6lyChyQTL2xALD6",
    ),
    (
        "loreX",
        "d8d40804-e1a1-4d22-b5f8-87eb11a5d4b4",
        "lore@codex.com",
        "lorex.jpeg",
        "$2a$14$Azk3S9Zkbt0LrV/uN8WVMuoFjl0MJxpuXp31ps6lyChyQTL2xALD6",
    ),
]

channels = [
    (
        1,
        "Go Concurrency Deep Dive",
        "dreamsofcode",
        "5f90bcae-b649-4b26-8f6c-1fa9217303f6",
        "doc.jpg",
        "go.png",
    ),
    (
        2,
        "CSS Advanced Techniques",
        "webDevSimplified",
        "94e52033-8ec2-413d-aa49-7a69d23f6836",
        "wds.png",
        "css.png",
    ),
    (
        3,
        "C# for Beginners",
        "kamil",
        "b920bfdc-7a44-445d-ba09-a665df542ac9",
        "kamil.png",
        "c#.png",
    ),
    (
        4,
        "Java Memory Management",
        "loreX",
        "d8d40804-e1a1-4d22-b5f8-87eb11a5d4b4",
        "lorex.jpeg",
        "java.png",
    ),
    (
        5,
        "JavaScript Event Loop Explained",
        "frontendfriends",
        "3da1bc7e-e9bd-4518-8151-6d5ff18d06c3",
        "kpowell.jpeg",
        "js.png",
    ),
    (
        6,
        "Rust Ownership Patterns",
        "dreamsofcode",
        "5f90bcae-b649-4b26-8f6c-1fa9217303f6",
        "doc.jpg",
        "rust.png",
    ),
    (
        7,
        "Python Data Science Hub",
        "webDevSimplified",
        "94e52033-8ec2-413d-aa49-7a69d23f6836",
        "wds.png",
        "python.png",
    ),
    (
        8,
        "Node.js Performance Tuning",
        "kamil",
        "b920bfdc-7a44-445d-ba09-a665df542ac9",
        "kamil.png",
        "node.png",
    ),
    (
        9,
        "React State Management",
        "loreX",
        "d8d40804-e1a1-4d22-b5f8-87eb11a5d4b4",
        "lorex.jpeg",
        "react.png",
    ),
    (
        10,
        "Vue Composition API",
        "frontendfriends",
        "3da1bc7e-e9bd-4518-8151-6d5ff18d06c3",
        "frotendfriends.png",
        "vue.png",
    ),
    (
        11,
        "Machine Learning in Practice",
        "dreamsofcode",
        "5f90bcae-b649-4b26-8f6c-1fa9217303f6",
        "doc.jpg",
        "ml.png",
    ),
    (
        12,
        "Docker and Containerization",
        "webDevSimplified",
        "94e52033-8ec2-413d-aa49-7a69d23f6836",
        "wds.png",
        "docker.png",
    ),
    (
        13,
        "Kubernetes for Devs",
        "kamil",
        "b920bfdc-7a44-445d-ba09-a665df542ac9",
        "kamil.png",
        "kubernetes.png",
    ),
    (
        14,
        "PostgreSQL Performance Tips",
        "loreX",
        "d8d40804-e1a1-4d22-b5f8-87eb11a5d4b4",
        "lorex.jpeg",
        "postgres.png",
    ),
    (
        15,
        "SQLite for Mobile Apps",
        "frontendfriends",
        "3da1bc7e-e9bd-4518-8151-6d5ff18d06c3",
        "kpowell.jpeg",
        "sqlite.png",
    ),
    (
        16,
        "Microservices Architecture",
        "dreamsofcode",
        "5f90bcae-b649-4b26-8f6c-1fa9217303f6",
        "doc.jpg",
        "microservices.png",
    ),
    (
        17,
        "API Security Fundamentals",
        "webDevSimplified",
        "94e52033-8ec2-413d-aa49-7a69d23f6836",
        "wds.png",
        "apisec.png",
    ),
    (
        18,
        "OAuth2 and Identity",
        "kamil",
        "b920bfdc-7a44-445d-ba09-a665df542ac9",
        "kamil.png",
        "oauth2.png",
    ),
    (
        19,
        "Web Accessibility Matters",
        "loreX",
        "d8d40804-e1a1-4d22-b5f8-87eb11a5d4b4",
        "lorex.jpeg",
        "webaccess.png",
    ),
    (
        20,
        "TypeScript Type Systems",
        "frontendfriends",
        "3da1bc7e-e9bd-4518-8151-6d5ff18d06c3",
        "kpowell.jpeg",
        "ts.png",
    ),
    (
        21,
        "DevOps CI/CD Strategies",
        "dreamsofcode",
        "5f90bcae-b649-4b26-8f6c-1fa9217303f6",
        "doc.jpg",
        "devops.png",
    ),
    (
        22,
        "Clean Code Principles",
        "webDevSimplified",
        "94e52033-8ec2-413d-aa49-7a69d23f6836",
        "wds.png",
        "cleancode.png",
    ),
    (
        23,
        "Design Patterns in Go",
        "kamil",
        "b920bfdc-7a44-445d-ba09-a665df542ac9",
        "kamil.png",
        "go.png",
    ),
    (
        24,
        "Functional Programming in JS",
        "loreX",
        "d8d40804-e1a1-4d22-b5f8-87eb11a5d4b4",
        "lorex.jpeg",
        "js.png",
    ),
    (
        25,
        "SvelteKit for Web Apps",
        "frontendfriends",
        "3da1bc7e-e9bd-4518-8151-6d5ff18d06c3",
        "kpowell.jpeg",
        "svelte.png",
    ),
]


def make_titles(channel_name):
    base = channel_name.split()[0]
    # 40 unique, topic-relevant titles
    return [
        f"{base} Fundamentals",
        f"Advanced {base} Techniques",
        f"{base} in the Real World",
        f"Common Pitfalls in {base}",
        f"Best Practices for {base}",
        f"{base} Performance Tuning",
        f"Debugging {base} Applications",
        f"{base} for Beginners",
        f"Expert Tips for {base}",
        f"{base} and Modern Development",
        f"{base} Security Essentials",
        f"{base} and Scalability",
        f"Testing in {base}",
        f"{base} and Cloud Integration",
        f"{base} in Production",
        f"{base} and Open Source",
        f"{base} Community Insights",
        f"{base} and Legacy Systems",
        f"{base} Migration Strategies",
        f"{base} and Automation",
        f"{base} Tooling Overview",
        f"{base} and Continuous Delivery",
        f"{base} and Monitoring",
        f"{base} and Logging",
        f"{base} and Error Handling",
        f"{base} and API Design",
        f"{base} and UI/UX",
        f"{base} and Mobile Development",
        f"{base} and Data Management",
        f"{base} and Security Audits",
        f"{base} and Performance Metrics",
        f"{base} and Testing Frameworks",
        f"{base} and Code Reviews",
        f"{base} and Refactoring",
        f"{base} and Version Control",
        f"{base} and Package Management",
        f"{base} and Dependency Injection",
        f"{base} and Design Patterns",
        f"{base} and Architecture",
        f"{base} and Future Trends",
        f"{base} and Career Growth",
    ]


def make_channel_desc(channel_name):
    ch_desc = f"{channel_name} aims to help you through all of its mysteries and help you to become entirely proficient in its use."
    return ch_desc


def make_content(channel_name):
    topic = channel_name.split()[0]
    paragraphs = [
        f"{topic} is a cornerstone of modern software development, offering developers a powerful set of tools and paradigms to build robust, scalable, and maintainable applications. In recent years, the {topic} ecosystem has evolved rapidly, introducing new features, libraries, and best practices that have transformed the way engineers approach problem-solving. Whether you're working on web applications, mobile apps, or backend systems, understanding the intricacies of {topic} can help you deliver high-quality solutions that meet the demands of today's technology landscape.",
        f"One of the key challenges in mastering {topic} is staying up-to-date with the latest advancements and community-driven innovations. Developers often face complex scenarios involving performance optimization, security, and integration with other technologies. By leveraging the collective knowledge shared in forums, open-source projects, and technical blogs, you can gain valuable insights into real-world use cases and practical solutions. As you deepen your expertise in {topic}, you'll discover patterns and strategies that not only improve your codebase but also enhance collaboration within your team.",
        f"In this channel, we explore a wide range of topics related to {topic}, from foundational concepts to advanced techniques. Our discussions cover everything from debugging and testing to deployment and monitoring, ensuring that you have the resources needed to tackle any challenge. Whether you're a seasoned professional or just starting your journey, the {channel_name} community is here to support your growth and help you achieve your goals in the ever-evolving world of technology.",
        f"Real-world applications of {topic} often require a blend of theoretical knowledge and hands-on experience. By participating in code reviews, contributing to open-source projects, and engaging in collaborative problem-solving, you can refine your skills and stay ahead of industry trends. The ability to adapt to new tools and methodologies is essential for long-term success, and continuous learning is the key to unlocking your full potential as a developer.",
        f"As you navigate the complexities of {topic}, remember that every challenge is an opportunity to learn and grow. Embrace the spirit of innovation, share your experiences with others, and never stop exploring the possibilities that {topic} has to offer. Together, we can build a vibrant community that empowers developers to create impactful solutions and shape the future of technology.",
    ]
    return "\n\n".join(paragraphs)


with open("insert_posts_and_channels.sql", "w") as f:
    user = 1
    for idx, (
        user_name,
        user_id,
        user_email,
        user_avatar,
        user_pass,
    ) in enumerate(users, 1):
        created = f"datetime('2025-07-28', '-{random.randint(1, 1000)} days')"
        owner_id_blob = user_id.replace("-", "")
        sql_u = (
            "INSERT INTO Users (ID, Username, EmailAddress, Avatar, Banner, Description, Usertype, Created, IsFlagged, SessionToken, CsrfToken, HashedPassword) VALUES (\n"
            "  CAST(X'{}' AS BLOB),\n"
            " '{}',\n "
            " '{}',\n "
            " '{}',\n "
            " 'default.png',\n "
            " '',\n "
            " 'user',\n "
            " {},\n "
            " 0,\n "
            " '',\n "
            " '',\n "
            " '{}'\n "
            ");\n"
        ).format(
            owner_id_blob,
            user_name,
            user_email,
            user_avatar,
            created,
            user_pass,
        )
        f.write(sql_u)
        user += 1

with open("insert_posts_and_channels.sql", "a") as f:
    for idx, (
        channel_id,
        channel_name,
        owner_username,
        owner_id,
        owner_avatar,
        channel_avatar,
    ) in enumerate(channels, 1):
        desc = make_channel_desc(channel_name)
        created = f"datetime('2025-07-28', '-{random.randint(1, 1000)} days')"
        owner_id_blob = owner_id.replace("-", "")
        sql_c = (
            "INSERT INTO Channels (OwnerID, Name, Description, Created, Avatar, Banner, Privacy, IsFlagged, IsMuted) VALUES (\n"
            "  CAST(X'{}' AS BLOB),\n"
            " '{}',\n "
            " '{}',\n "
            " {},\n "
            " '{}',\n "
            " 'default.png',\n"
            " 0,\n "
            " 0,\n "
            " 0\n "
            ");\n"
        ).format(
            owner_id_blob,
            channel_name,
            desc,
            created,
            channel_avatar,
        )
        f.write(sql_c)

with open("insert_posts_and_channels.sql", "a") as f:
    post_id = 1
    for idx, (
        channel_id,
        channel_name,
        owner_username,
        owner_id,
        owner_avatar,
        channel_avatar,
    ) in enumerate(channels, 1):
        titles = make_titles(channel_name)
        content = make_content(channel_name)
        for i, title in enumerate(titles, 1):
            full_title = f"{title}"
            created = f"datetime('2025-07-28', '-{random.randint(1, 1000)} days')"
            author_id_blob = owner_id.replace("-", "")
            sql_p = (
                "INSERT INTO Posts (Title, Content, Images, Created, IsCommentable, IsFlagged, Author, AuthorID, AuthorAvatar) VALUES (\n"
                "  '{}',\n"
                "  '{}',\n"
                "  'noimage',\n"
                "  {},\n"
                "  1,\n"
                "  0,\n"
                "  '{}',\n"
                "  CAST(X'{}' AS BLOB),\n"
                "  '{}'\n"
                ");\n"
            ).format(
                full_title.replace("'", "''"),
                content.replace("'", "''"),
                created,
                owner_username,
                author_id_blob,
                owner_avatar,
            )
            f.write(sql_p)
            sql_c = (
                "INSERT INTO PostChannels (ChannelID, PostID, Created) VALUES (\n"
                "  {},\n"
                "  {},\n"
                "  {}\n"
                ");\n"
            ).format(
                channel_id,
                post_id,
                created,
            )
            f.write(sql_c)
            post_id += 1
