import {Component, Input, OnInit} from '@angular/core';

@Component({
  selector: 'app-linkcard',
  templateUrl: './linkcard.component.html',
  styleUrls: ['./linkcard.component.css']
})
export class LinkcardComponent implements OnInit {

  @Input() image: string;
  @Input() text: string;
  @Input() linkTo: string;

  constructor() {
  }

  ngOnInit() {
  }

}
