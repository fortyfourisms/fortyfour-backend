DELETE FROM risiko
WHERE kode IN (
    'ip-theft',
    'data-breach',
    'kerusakan-perangkat',
    'kehilangan-peralatan',
    'human-error',
    'third-party-breach',
    'management-lack-tech',
    'malware-attack',
    'ddos-attack',
    'phishing',
    'zero-day',
    'ransomware',
    'brute-force'
);
