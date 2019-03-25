import {Component, OnInit} from '@angular/core';
import {PluginsService} from 'src/app/shared/services/plugins.service';
import {Plugin} from 'src/app/models/plugin';
import {CeRegisterService} from "../shared/services/ce-register.service";
import {CELoadedEvent, EventType} from "../models/events";
import {Observable} from "rxjs";

@Component({
  selector: 'app-settings',
  templateUrl: './settings.component.html',
  styleUrls: ['./settings.component.scss'],
})
export class SettingsComponent {

  private plugins :Observable<CELoadedEvent>;

  constructor(
    private pluginsService: PluginsService,
    private ceRegisterService: CeRegisterService,
  ){}

  ngAfterViewInit() {
    this.pluginsService.getAll().map((plugin: Plugin) => {
      console.log(plugin);
      this.plugins = this.ceRegisterService.registerCe(plugin.checkComponentEntryPoint);

      this.plugins.subscribe((event: CELoadedEvent) => {
        const pluginCustomEl: HTMLElement = document.createElement(event.selector);

        pluginCustomEl.addEventListener(EventType.ACTION_SUBMIT_EVENT, msg => console.debug('plugin actionSubmit says: ', msg));

        const pluginsContainer = document.querySelector('app-settings');
        pluginsContainer.appendChild(pluginCustomEl);

        setTimeout(function () {
          pluginCustomEl.setAttribute('checkResult', 'init');
          console.log('checkResult was sent');
        }, 2000);
        });

    });



  }
}
