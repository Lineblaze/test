import http from 'k6/http';
import { check, sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export let options = {
    stages: [
        { duration: '1m', target: 500 },
        { duration: '2m', target: 1000 },
        { duration: '1m', target: 0 },
    ],
    thresholds: {
        http_req_duration: ['p(95)<50'],
        http_req_failed: ['rate<0.0001'],
    },
};

export default function () {
    const url = 'http://localhost:8080/api/auth';

    const username = `user_${randomString(8)}`;
    const password = `pass_${randomString(12)}`;

    const payload = JSON.stringify({ username, password });

    const params = { headers: { 'Content-Type': 'application/json' } };

    const res = http.post(url, payload, params);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'response time is acceptable': (r) => r.timings.duration < 50,
    });

    sleep(1);
}
