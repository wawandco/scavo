-- Seed data: 4 teams, 12 players (3 per team) + 1 admin, 5 Central Park hunts

INSERT INTO teams (id, name, created_at, updated_at) VALUES
(1, 'The Park Rangers',      datetime('now'), datetime('now')),
(2, 'Central Park Striders', datetime('now'), datetime('now')),
(3, 'Bethesda Brigade',      datetime('now'), datetime('now')),
(4, 'The Ramblers',          datetime('now'), datetime('now'));

INSERT INTO members (id, name, personal_id, team_id, is_admin, created_at, updated_at) VALUES
(1,  'Alex Chen',         'alex',      1, 0, datetime('now'), datetime('now')),
(2,  'Maya Patel',        'maya',      1, 0, datetime('now'), datetime('now')),
(3,  'James Wilson',      'james',     1, 0, datetime('now'), datetime('now')),
(4,  'Sofia Ramirez',     'sofia',     2, 0, datetime('now'), datetime('now')),
(5,  'Liam O''Brien',     'liam',      2, 0, datetime('now'), datetime('now')),
(6,  'Emma Tanaka',       'emma',      2, 0, datetime('now'), datetime('now')),
(7,  'Noah Kim',          'noah',      3, 0, datetime('now'), datetime('now')),
(8,  'Olivia Davis',      'olivia',    3, 0, datetime('now'), datetime('now')),
(9,  'Ethan Brooks',      'ethan',     3, 0, datetime('now'), datetime('now')),
(10, 'Isabella Garcia',   'isabella',  4, 0, datetime('now'), datetime('now')),
(11, 'Lucas Anderson',    'lucas',     4, 0, datetime('now'), datetime('now')),
(12, 'Charlotte Martin',  'charlotte', 4, 0, datetime('now'), datetime('now')),
(13, 'Admin User',        'admin',     NULL, 1, datetime('now'), datetime('now'));

INSERT INTO hunts (id, title, description, points, created_at, updated_at) VALUES
(1, 'The Imagine Mosaic',
    'Snap a photo of the iconic "Imagine" mosaic in Strawberry Fields, a tribute to John Lennon.',
    100, datetime('now'), datetime('now')),
(2, 'Bow Bridge Selfie',
    'Take a selfie on the Bow Bridge with The Lake visible behind you — the most romantic spot in the park.',
    200, datetime('now'), datetime('now')),
(3, 'Belvedere Castle View',
    'Climb to the top of Belvedere Castle and capture a photo of the Great Lawn and the Manhattan skyline.',
    300, datetime('now'), datetime('now')),
(4, 'Bethesda Fountain & Angel',
    'Photograph the Angel of the Waters statue at Bethesda Fountain, the heart of the park.',
    400, datetime('now'), datetime('now')),
(5, 'The Mall Literary Walk',
    'Find the statue of Hans Christian Andersen and take a photo reading a story with him among the elm trees.',
    500, datetime('now'), datetime('now'));