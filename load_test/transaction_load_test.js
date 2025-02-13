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

const items = [
    { id: 't-shirt', price: 80 },
    { id: 'cup', price: 20 },
    { id: 'book', price: 50 },
    { id: 'pen', price: 10 },
    { id: 'powerbank', price: 200 },
    { id: 'hoody', price: 300 },
    { id: 'umbrella', price: 200 },
    { id: 'socks', price: 10 },
    { id: 'wallet', price: 50 },
    { id: 'pink-hoody', price: 500 },
];

const users = Array.from({ length: 10 }, (_, i) => `user_${i + 1}`);

function getAuthToken() {
    const authUrl = 'http://localhost:8080/api/auth';
    const username = `user_${randomString(8)}`;
    const password = `pass_${randomString(12)}`;
    const payload = JSON.stringify({ username, password });

    const params = { headers: { 'Content-Type': 'application/json' } };
    const res = http.post(authUrl, payload, params);

    if (res.status === 200) {
        const jsonResponse = JSON.parse(res.body);
        return jsonResponse.token;
    }
    return null;
}

export default function () {
    const token = getAuthToken();
    if (!token) {
        console.error('Failed to authenticate');
        return;
    }

    const headers = {
        'Content-Type': 'application/json',
        'Authorization': `Bearer ${token}`,
    };

    const randomItem = items[Math.floor(Math.random() * items.length)];
    const buyRes = http.get(`http://localhost:8080/transactions/buy/${randomItem.id}`, { headers });
    check(buyRes, {
        'Buy request successful': (r) => r.status === 200,
        'Response time < 50ms': (r) => r.timings.duration < 50,
    });

    const randomUser = users[Math.floor(Math.random() * users.length)];
    const sendPayload = JSON.stringify({ recipient: randomUser, amount: Math.floor(Math.random() * 100 + 1) });
    const sendRes = http.post('http://localhost:8080/transactions/send', sendPayload, { headers });
    check(sendRes, {
        'Send request successful': (r) => r.status === 200,
        'Response time < 50ms': (r) => r.timings.duration < 50,
    });

    const infoRes = http.get('http://localhost:8080/transactions/info', { headers });
    check(infoRes, {
        'Info request successful': (r) => r.status === 200,
        'Response time < 50ms': (r) => r.timings.duration < 50,
    });

    sleep(1);
}
