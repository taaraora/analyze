import { Component, OnInit } from '@angular/core';
import { PluginsService }    from "src/app/plugins/plugins.service";

@Component({
  selector: 'app-plugins',
  templateUrl: './plugins.component.html',
  styleUrls: ['./plugins.component.scss']
})
export class PluginsComponent implements OnInit {
  plugins: any;

  constructor(
    private pluginsService: PluginsService
  ) { }

  ngOnInit() {
    this.plugins = this.pluginsService.all();
  }

}
