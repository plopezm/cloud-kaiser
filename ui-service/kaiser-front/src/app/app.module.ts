import {BrowserModule} from '@angular/platform-browser';
import {NgModule} from '@angular/core';
import {BrowserAnimationsModule} from '@angular/platform-browser/animations';
import {MatToolbarModule} from '@angular/material/toolbar';
import {MatSidenavModule} from '@angular/material/sidenav';
import {MatButtonModule} from '@angular/material/button';
import {MatIconModule} from '@angular/material/icon';
import {MatCardModule} from '@angular/material/card';

import {AppComponent} from './app.component';
import {HomeComponent} from './containers/home/home.component';
import {AppRoutingModule} from './app-routing.module';
import {MenubarComponent} from './containers/menubar/menubar.component';
import {LinkcardComponent} from './components/cards/linkcard/linkcard.component';
import {FlexLayoutModule} from "@angular/flex-layout";
import {HttpClientModule} from "@angular/common/http";

@NgModule({
  imports: [
    HttpClientModule,
    BrowserModule,
    AppRoutingModule,
    BrowserAnimationsModule,
    FlexLayoutModule,
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
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule {
}
