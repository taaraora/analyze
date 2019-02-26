import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { IPlugin }    from 'src/app/models/models';
import { Observable } from 'rxjs';

@Injectable()
export class PluginsService {

  constructor(
    private http: HttpClient,
  ) {
  }

  all(): Observable<IPlugin[]> {

    return this.http.get<IPlugin[]>('/api/v1/plugins');
  }

  loadExternal(path: string): void {
    const header = document.querySelector('head');
    // wtf
    header.addEventListener('loadingNotifier', (msg: CustomEvent) => {
      //get event that it was loaded
      const pluginCustomEl: HTMLElement = document.createElement(msg.detail.selector);

      pluginCustomEl.addEventListener('actionSubmit', msg => console.debug('plugin actionSubmit says: ', msg));

      const pluginsContainer = document.querySelector('app-plugins');
      pluginsContainer.appendChild(pluginCustomEl);

      setTimeout(function () {
        pluginCustomEl.setAttribute('checkResult', 'init');
        console.log('checkResult was sent');
      }, 2000);
    });

    const script = document.createElement('script');
    // wtf
    script.src = path;
    header.appendChild(script);
  }
}
