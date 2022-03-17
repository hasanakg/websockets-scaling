const os = require('os');
const ifaces = os.networkInterfaces();

const privateIp = (() => {
  return Object.values(ifaces).flat().find(val => {
    return (val.family == 'IPv4' && val.internal == false);
  }).address;
})();
let connectCounter = 0;
const randomOffset = Math.floor(Math.random() * 10);
const intervalOffset = (30+randomOffset) * Math.pow(10,3);

// WebSocket Server
const socketPort = 5000;
const socketServer = require('http').createServer();
const io = require('socket.io')(socketServer, {
  path: '/'
});
// Redis Adapter
const redisAdapter = require('socket.io-redis');
const redisHost = process.env.REDIS_HOST || 'localhost'; // 'redis-master.default.svc';
io.adapter(redisAdapter({ host: redisHost, port: 6379 }));

let clientId = undefined;

// Handlers
io.on('connection', client => {
  console.log('New incoming Connection from', client.id);
  console.log( client.conn.server.clientsCount + " users connected" );
  client.on('test000', function(message) {
    console.log('Message from the client:',client.id,'->',message);
    console.log( client.conn.server.clientsCount + " users connected" );
    clientId = client.id;
    io.to(client.id).emit("okok", `Your ID is ${client.id} and I'm ${privateIp}`);
    // setTimeout(() => {
    //   io.to(client.id).emit("disconnect");
    // }, 5000);
  })
  // io.on('connect', function() { connectCounter++; console.log(`Connected clients up: ${connectCounter}`); });
  // io.on('disconnect', function() { connectCounter--; console.log(`Connected clients down: ${connectCounter}`); });

});
setInterval(() => {
  let log0 = `I am the host: ${privateIp}. I am healty.`;
  console.log(log0);
  io.emit("okok", log0);
}, intervalOffset);

// Web Socket listen
socketServer.listen(socketPort);