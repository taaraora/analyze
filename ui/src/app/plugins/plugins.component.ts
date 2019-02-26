import { Component, OnInit } from '@angular/core';
import { PluginsService }    from 'src/app/plugins/plugins.service';
import { tap }               from 'rxjs/operators';
import { IPlugin }           from 'src/app/models/models';

@Component({
  selector: 'app-plugins',
  templateUrl: './plugins.component.html',
  styleUrls: ['./plugins.component.scss'],
})
export class PluginsComponent implements OnInit {
  plugins$: any;


  constructor(
    private pluginsService: PluginsService,
  ) {
  }

  ngOnInit() {
    console.log('init');
    this.plugins$ = this.pluginsService.all();
  }

  ngAfterViewInit() {
    this.plugins$.pipe(
      tap((plugin: IPlugin) => this.pluginsService.loadExternal(plugin.settingsComponentEntryPoint)),
    ).subscribe();
  }
}
