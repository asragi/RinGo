CREATE DATABASE IF NOT EXISTS ringo;
CREATE TABLE IF NOT EXISTS ringo.users(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` varchar(40) NOT NULL,
    `name` varchar(40) NOT NULL,
    `fund` int(20) NOT NULL,
    `max_stamina` int(11) NOT NULL,
    `stamina_recover_time`DATETIME NOT NULL,
    `hashed_password` varchar(40)
    PRIMARY KEY (`id`),
    INDEX `user_id_index` (`user_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.item_masters(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `item_id` varchar(40) NOT NULL,
    `display_name` varchar(40) NOT NULL,
    `description` varchar(40) NOT NULL,
    `price` int(20) NOT NULL,
    `max_stock` int(10) NOT NULL,
    `sale_freq` float(8,4) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `item_id_index` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.item_storages(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` varchar(40) NOT NULL,
    `item_id` varchar(40) NOT NULL,
    `stock` int(10) NOT NULL,
    `is_known` bool NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `user_id_item_id_index` (`user_id`, `item_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    FOREIGN KEY (`item_id`) REFERENCES `item_masters` (`item_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.skill_masters(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `skill_id` varchar(40) NOT NULL,
    `display_name` varchar(40) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `skill_id_index` (`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.user_skills(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `skill_id` varchar(40) NOT NULL,
    `user_id` varchar(40) NOT NULL,
    `skill_exp` int(20) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `user_id_skill_id_index` (`user_id`, `skill_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill_masters` (`skill_id`),
    CONSTRAINT user_skill_pair UNIQUE (`user_id`, `skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.explore_masters(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `explore_id` varchar(40) NOT NULL,
    `display_name` varchar(40) NOT NULL,
    `description` varchar(40) NOT NULL,
    `consuming_stamina` int(10) NOT NULL,
    `required_payment` int(10) NOT NULL,
    `stamina_reducible_rate` float(6,5) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `explore_id_index` (`explore_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.user_explore_data(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `user_id` varchar(40) NOT NULL,
    `explore_id` varchar(40) NOT NULL,
    `is_known` bool NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `explore_id_index` (`explore_id`),
    FOREIGN KEY (`user_id`) REFERENCES `users` (`user_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.skill_growth_data(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `explore_id` varchar(40) NOT NULL,
    `skill_id` varchar(40) NOT NULL,
    `gaining_point` int(20) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `user_id_skill_id_index` (`explore_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill_masters` (`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.stage_masters(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `stage_id` varchar(40) NOT NULL,
    `display_name` varchar(40) NOT NULL,
    `description` varchar(40) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `stage_id_index` (`stage_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.stage_explore_relations(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `stage_id` varchar(40) NOT NULL,
    `explore_id` varchar(40) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `stage_id_index` (`stage_id`),
    FOREIGN KEY (`stage_id`) REFERENCES `stage_masters` (`stage_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.item_explore_relations(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `item_id` varchar(40) NOT NULL,
    `explore_id` varchar(40) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `item_id_index` (`item_id`, `explore_id`),
    FOREIGN KEY (`item_id`) REFERENCES `item_masters` (`item_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.earning_items(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `explore_id` varchar(40) NOT NULL,
    `item_id` varchar(40) NOT NULL,
    `min_count` int(10) NOT NULL,
    `max_count` int(10) NOT NULL,
    `probability` float(6,5) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `explore_id_index` (`explore_id`),
    FOREIGN KEY (`item_id`) REFERENCES `item_masters` (`item_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.consuming_items(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `explore_id` varchar(40) NOT NULL,
    `item_id` varchar(40) NOT NULL,
    `max_count` int(10) NOT NULL,
    `consumption_prob` float(6,5) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `explore_id_index` (`explore_id`),
    FOREIGN KEY (`item_id`) REFERENCES `item_masters` (`item_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.required_skills(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `explore_id` varchar(40) NOT NULL,
    `skill_id` varchar(40) NOT NULL,
    `skill_lv` int(4) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `explore_id_index` (`explore_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill_masters` (`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;

CREATE TABLE IF NOT EXISTS ringo.stamina_reduction_skills(
    `id` int(11) NOT NULL AUTO_INCREMENT,
    `explore_id` varchar(40) NOT NULL,
    `skill_id` varchar(40) NOT NULL,
    PRIMARY KEY (`id`),
    INDEX `explore_id_index` (`explore_id`, `skill_id`),
    FOREIGN KEY (`explore_id`) REFERENCES `explore_masters` (`explore_id`),
    FOREIGN KEY (`skill_id`) REFERENCES `skill_masters` (`skill_id`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
