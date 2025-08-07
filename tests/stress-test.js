import http from 'k6/http';

export const options = {
    vus: 5000,
    duration: '30s',
};

export default () => {
  http.get('http://k8s.orb.local/healthz');
};