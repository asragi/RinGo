USE `ringo`;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/item-master.csv'
    INTO TABLE `item_masters`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/skill-master.csv'
    INTO TABLE `skill_masters`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/explore-master.csv'
    INTO TABLE `explore_masters`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/skill-growth.csv'
    INTO TABLE `skill_growth_data`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/stage-master.csv'
    INTO TABLE `stage_masters`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/stage-explore-relations.csv'
    INTO TABLE `stage_explore_relations`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/item-explore-relations.csv'
    INTO TABLE `item_explore_relations`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/earning-items.csv'
    INTO TABLE `earning_items`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/consuming-items.csv'
    INTO TABLE `consuming_items`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/required-skills.csv'
    INTO TABLE `required_skills`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
LOAD DATA LOCAL INFILE '/docker-entrypoint-initdb.d/data/reduction-stamina.csv'
    INTO TABLE `stamina_reduction_skills`
    FIELDS TERMINATED BY ',' ESCAPED BY '"'
    LINES TERMINATED BY '\r\n'
    IGNORE 1 ROWS;
