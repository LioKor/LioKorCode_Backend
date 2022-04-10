(function () {
  'use strict';

  function SocketConnection (url) {
    this.url = url;
    this.ws = null;

    var ws = this.ws = new WebSocket(this.url);
    var self = this;

    ws.onopen = function (evt) {
      self.emit('open', evt);
    };

    ws.onclose = function (evt) {
      self.emit('close', evt);
    };

    ws.onmessage = function (evt) {
      console.log(JSON.parse(evt.data))
      var m = JSON.parse(evt.data);

      if (m && m.e) {
        console.log(JSON.parse(evt.data))
        self.emit(m.e, m.d, evt);
      }
    };
  }

  SocketConnection.prototype = new EventEmitter;

  SocketConnection.prototype.send = function(eventName, data) {
    console.log(JSON.stringify({e: eventName, d: data}));
    this.ws.send(JSON.stringify({e: eventName, d: data}));
  };

  window.SocketConnection = SocketConnection;
}());
