INSERT INTO transcripts (id, `text`, `role`, simulation_id) VALUES
    -- FIRST SIMULATION
    (1, 'Malta Ground, RYR4735, Stand 4, request taxi.', 'aircraft', 1),
    (2, 'RYR4735, taxi to holding point C via Victor, Uniform, Charlie, contact Tower on 135.105.', 'tower', 1),
    (3, 'Malta Tower, RYR4735, holding at C, ready for departure.', 'aircraft', 1),
    (4, 'RYR4735, roger, line-up runway 31, cleared for takeoff.', 'tower', 1),
    (5, 'Cleared for takeoff runway 31, RYR4735.', 'aircraft', 1),
    (6, 'Luqa Radar, RYR4735, passing 2000 feet, climbing 5000ft, on the GZO3D SID.', 'aircraft', 1),

    -- SECOND SIMULATION
    (7, 'Malta Radar, KMM514, maintaining FL340, inbound ARLOS.', 'aircraft', 2),
    (8, 'KMM514, radio contact, maintain FL340, proceed direct EVIRA.', 'tower', 2),
    (16, 'Direct EVIRA, KMM514.', 'aircraft', 2),

    -- THIRD SIMULATION
    (9, 'Malta Approach, KMM102, descending to FL150, approaching EKOLA.', 'aircraft', 3),
    (10, 'KMM102, radar contact. Proceed direct KEKOR, descend FL70, expect ILS approach runway 31.', 'tower', 3),
    (11, 'Direct KEKOR, descending to FL70, expecting ILS runway 31, KMM102.', 'aircraft', 3),
    (12, 'KMM102, descend to altitude 3,000 feet, via KEKOR, cleared ILS approach runway 31.', 'tower', 3),
    (13, 'Descending to 3,000 feet, via KEKOR, cleared ILS approach runway 31, KMM102.', 'aircraft', 3),
    (14, 'KMM102, continue approach with Luqa Tower on 135.105.', 'tower', 3),
    (15, 'Tower on 135.105, KMM102.', 'aircraft', 3);






