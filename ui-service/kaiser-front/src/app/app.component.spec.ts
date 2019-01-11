import {async, TestBed} from '@angular/core/testing';
import {AppComponent} from './app.component';
import {MatButtonModule, MatCardModule, MatIconModule, MatSidenavModule, MatToolbarModule} from "@angular/material";
import {AppRoutingModule} from "./app-routing.module";
import {HomeComponent} from "./containers/home/home.component";
import {MenubarComponent} from "./containers/menubar/menubar.component";
import {LinkcardComponent} from "./components/cards/linkcard/linkcard.component";
import {BrowserAnimationsModule} from "@angular/platform-browser/animations";

describe('AppComponent', () => {
  beforeEach(async(() => {
    TestBed.configureTestingModule({
      imports: [
        BrowserAnimationsModule,
        AppRoutingModule,
        MatToolbarModule,
        MatSidenavModule,
        MatButtonModule,
        MatIconModule,
        MatCardModule
      ],
      declarations: [
        AppComponent,
        HomeComponent,
        MenubarComponent,
        LinkcardComponent
      ],

    }).compileComponents();
  }));

  it('should create the app', () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app).toBeTruthy();
  });

  it(`should have as title 'kaiser-front'`, () => {
    const fixture = TestBed.createComponent(AppComponent);
    const app = fixture.debugElement.componentInstance;
    expect(app.title).toEqual('kaiser-front');
  });
});
