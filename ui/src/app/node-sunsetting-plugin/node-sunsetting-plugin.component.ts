import {ChangeDetectorRef, Component, Input, OnInit} from '@angular/core';

@Component({
  // selector: 'app-node-sunsetting-plugin',
  templateUrl: './node-sunsetting-plugin.component.html',
  styleUrls: ['./node-sunsetting-plugin.component.scss']
})
export class NodeSunsettingPluginComponent implements OnInit {

  @Input() a: number;
  @Input() b: number;
  @Input() c: number;

  constructor(private cd: ChangeDetectorRef) { }

  ngOnInit() {
  }

  more(): void {
    this.a = Math.round(Math.random() * 100);
    this.b = Math.round(Math.random() * 100);
    this.c = Math.round(Math.random() * 100);

    this.cd.markForCheck();
  }

}
