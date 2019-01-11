import {async, ComponentFixture, TestBed} from '@angular/core/testing';

import {MenubarComponent} from './menubar.component';
import {MatButtonModule, MatCardModule, MatIconModule, MatSidenavModule, MatToolbarModule} from "@angular/material";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

describe('MenubarComponent', () => {
  let component: MenubarComponent;
  let fixture: ComponentFixture<MenubarComponent>;

  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        BrowserAnimationsModule,
        MatToolbarModule,
        MatSidenavModule,
        MatButtonModule,
        MatIconModule,
        MatCardModule
      ],
      declarations: [MenubarComponent]
    })
      .compileComponents();
  }));

  beforeEach(() => {
    fixture = TestBed.createComponent(MenubarComponent);
    component = fixture.componentInstance;
    fixture.detectChanges();
  });

  it('should create', () => {
    expect(component).toBeTruthy();
  });
});
