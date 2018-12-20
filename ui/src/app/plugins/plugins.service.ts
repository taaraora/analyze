import {Injectable} from '@angular/core';
import {HttpClient} from "@angular/common/http";

@Injectable()
export class PluginsService {

  constructor(
    private http: HttpClient
  ) { }

  all() {
    return this.http.get('/api/v1/plugin');
  }

  loadExternal(): void {
    console.log('test1');
    const header = document.querySelector('head');
    header.addEventListener('loadingNotifier', (msg: CustomEvent) => {
      //get event that it was loaded
      const element: HTMLElement = document.createElement(msg.detail.selector);
      element.addEventListener('actionSubmit', msg => console.debug('plugin actionSubmit says: ',msg));
      const pluginsContainer = document.querySelector('app-plugins');
      pluginsContainer.appendChild(element);
      console.log('test2');

      setTimeout(function () {
        element.setAttribute('checkResult', 'init');
        console.log('checkResult was sent')
      }, 2000);
    });

    const script = document.createElement('script');
    script.src = 'http://127.0.0.1:8080/main.js'
    header.appendChild(script);

  }
}
