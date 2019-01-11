import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {LinkcardComponent} from './linkcard.component';
import {MatButtonModule, MatCardModule, MatIconModule, MatSidenavModule, MatToolbarModule} from "@angular/material";

describe('LinkcardComponent', () => {
  let component: LinkcardComponent;
  let fixture: ComponentFixture<LinkcardComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        MatToolbarModule,
        MatSidenavModule,
        MatButtonModule,
        MatIconModule,
        MatCardModule
      ],
      declarations: [LinkcardComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(LinkcardComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
