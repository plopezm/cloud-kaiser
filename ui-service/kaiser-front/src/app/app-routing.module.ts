import {NgModule} from '@angular/core';
import {RouterModule, Routes} from '@angular/router';
import {HomeComponent} from './containers/home/home.component';
import {ManagementComponent} from "./containers/management/management.component";

const routes: Routes = [
  {path: '', component: HomeComponent},
  {path: 'management', component: ManagementComponent},
  {path: '**', redirectTo: '', pathMatch: 'full'},
];

@NgModule({
  declarations: [],
  imports: [
    RouterModule.forRoot(routes, {
      useHash: true
    })
  ],
  exports: [
    RouterModule
  ]
})
export class AppRoutingModule {
}
