import ws from 'k6/ws';
import { check } from 'k6';

export const options = {
    vus: 100000,
    duration: '1s',
  };

export default function () {
  const url = 'ws://localhost:8080/goapp/ws';
  const params = {};

  const res = ws.connect(url, params, function (socket) {
    socket.on('open', () => console.log('connected'));
    socket.on('message', (data) => console.log('Message received: ', data));
    socket.on('close', () => console.log('disconnected'));
    socket.send({});
  });

  check(res, { 'status is 101': (r) => r && r.status === 101 });
}