import { check } from 'k6';
import http from 'k6/http';
export let options = {
  duration:'10s',
  vus: 100
};
export default function() {
  let res = http.get('http://127.0.0.1:8080/ping');
  check(res, {
    'is status 200':(r) => r.status === 200
  });
};
