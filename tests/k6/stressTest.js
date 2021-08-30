import http from 'k6/http';
import { sleep } from 'k6';

export let options = {
  insecureSkipTLSVerify: true,
  noConnectionReuse: false,
  stages: [
    { duration: '2m', target: 100 },
    { duration: '5m', target: 100 },
    { duration: '2m', target: 200 },
    { duration: '5m', target: 200 },
    { duration: '2m', target: 400 },
    { duration: '5m', target: 400 },
    { duration: '2m', target: 600 },
    { duration: '5m', target: 600 },
  ],
};

export default () => {
  const url = 'http://localhost/events';
  const payload = JSON.stringify([
    {
      data: { id: 34, source: 2, type: 2, x: 1195, y: 524 },
      timestamp: 1627683336586,
      type: 3,
    },
  ]);

  const params = {
    headers: {
      'Content-Type': 'application/json',
      'Session-ID': '91946562-9973-4a2e-a72b-d1c7e670b034',
      'User-ID': '4d3f01e9-3718-4002-9331-85215dbb0e5e',
      'Chapter-ID': '5f9f01e9-3718-3023-9331-85216dbb0e5e',
      'X-Rrweb': 'true',
    },
  };

  http.post(url, payload, params);

  sleep(1);
};
