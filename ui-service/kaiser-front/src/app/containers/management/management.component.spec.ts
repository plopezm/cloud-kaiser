import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {ManagementComponent} from './management.component';
import {MatPaginatorModule, MatTableModule} from "@angular/material";
import {HttpClientModule} from "@angular/common/http";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

describe('ManagementComponent', () => {
  let component: ManagementComponent;
  let fixture: ComponentFixture<ManagementComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        BrowserAnimationsModule,
        HttpClientModule,
        MatTableModule,
        MatPaginatorModule
      ],
      declarations: [ManagementComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(ManagementComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
