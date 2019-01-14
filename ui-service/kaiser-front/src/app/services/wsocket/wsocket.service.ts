import {EventEmitter, Injectable, OnDestroy} from '@angular/core';
import {ConfigService} from "../config/config.service";

@Injectable({
  providedIn: 'root'
})
export class WebSocketService implements OnDestroy {

  private socket: WebSocket;
  private channels: any;

  constructor(private configService: ConfigService) {
    this.channels = {};
    this.socket = new WebSocket(configService.getWSEndpoint());
    /*
    this.socket.onopen = event => {
      this.listener.emit({"type": "open", "data": event});
    };
    this.socket.onclose = event => {
      this.listener.emit({"type": "close", "data": event});
    };
    */
    this.socket.onmessage = event => {
      //this.listener.emit({"type": "message", "data": JSON.parse(event.data)});
      let content = JSON.parse(event.data);
      let listener = this.channels[content.type];
      if (listener) {
        listener.emit(content);
      }
    }
  }

  ngOnDestroy(): void {
    debugger;
    this.socket.close();
  }

  public send(data: string) {
    this.socket.send(data);
  }

  public close() {
    this.socket.close();
  }

  public onEvent(channel: string) {
    let listener = this.channels[channel];
    if (!listener) {
      listener = new EventEmitter();
      this.channels[channel] = listener;
    }
    return listener
  }

}
