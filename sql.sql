CREATE TABLE `stats` (
    `int1` INT,
    `int2` INT,
    `limit` INT,
    `str1` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
    `str2` VARCHAR(255) COLLATE utf8mb4_unicode_ci,
    `hits` INT,
    PRIMARY KEY (`int1`,`int2`,`limit`,`str1`,`str2`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;