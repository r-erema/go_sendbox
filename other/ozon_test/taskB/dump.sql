CREATE TABLE IF NOT EXISTS goods (
  id    INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
  name  TEXT
);
CREATE TABLE  IF NOT EXISTS tags (
  id    INTEGER PRIMARY KEY AUTOINCREMENT UNIQUE,
  name  TEXT
);
CREATE TABLE IF NOT EXISTS tags_goods (
  tag_id INTEGER,
  goods_id INTEGER,
  UNIQUE (tag_id, goods_id)
);

INSERT INTO goods (`id`, `name`) VALUES
                                        (1, 'good_1'), (2, 'good_2'), (3, 'good_3'),
                                        (4, 'good_4'), (5, 'good_5'), (6, 'good_6'),
                                        (7, 'good_7'), (8, 'good_8'), (9, 'good_9'), (10, 'good_10');

INSERT INTO tags (`id`, `name`) VALUES (1, 'tag_1'), (2, 'tag_2'), (3, 'tag_3');
INSERT INTO tags_goods (`tag_id`, `goods_id`) VALUES (1, 2),
                                                    (1, 4),
                                                    (2, 1),
                                                    (3, 5),
                                                    (2, 5),
                                                    (3, 8),
                                                    (1, 8),
                                                    (2, 4),
                                                    (2, 2),
                                                    (3, 2),
                                                    (1, 9),
                                                    (2, 9),
                                                    (3, 4);
