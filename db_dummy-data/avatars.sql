UPDATE Posts 
SET AuthorAvatar = 
  CASE 
    WHEN AuthorID = 1 THEN 'EehNbqo0-KWx4RN-4R6V1Q==.jpeg'
    WHEN AuthorID = 2 THEN 'K1xcM1YYVT6xkHoEUICS4g==.png'
    WHEN AuthorID = 3 THEN 'noimage'
    WHEN AuthorID = 4 THEN 'A5JDosg0v9bCU8KOwjpmcw==.png'
    WHEN AuthorID = 5 THEN 'pKIil4ZLTuk2K3CKWBN9uQ==.jpeg'
  END;
