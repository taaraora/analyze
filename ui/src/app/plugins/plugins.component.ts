import { Component, OnInit } from '@angular/core';
import { PluginsService }    from "src/app/plugins/plugins.service";

@Component({
  selector: 'app-plugins',
  templateUrl: './plugins.component.html',
  styleUrls: ['./plugins.component.scss']
})
export class PluginsComponent implements OnInit {
  pluginsStatuses: any;


  constructor(
    private pluginsService: PluginsService
  ) { }

  ngOnInit() {
    this.pluginsStatuses = this.pluginsService.all();
  }

  ngAfterViewInit(){
    this.pluginsService.loadExternal();
  }

}
