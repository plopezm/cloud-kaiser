import { Injectable } from '@angular/core';
import {HttpClient} from "@angular/common/http";
import {TaskAdapter} from "../../core/models/task/task.adapter";
import {map} from "rxjs/operators";
import {Observable} from "rxjs";
import {Task} from "../../core/models/task/task.model";
import {ConfigService} from "../config/config.service";

@Injectable({
  providedIn: 'root'
})
export class TaskService {

  constructor(private http: HttpClient, private configService: ConfigService, private taskAdapter: TaskAdapter) { }

  getAllTasks(offset: number, limit: number): Observable<Task[]> {
    return this.http.get(this.configService.getTasksURL()).pipe(
        map((data: any[]) => data.map((item: any) => this.taskAdapter.adapt(item)))
    );
  }

  getTasksByName(name: string): Observable<Task[]> {
    return this.http.get(this.configService.getTasksByNameURL(name)).pipe(
      map((data: any[]) => data.map((item: any) => this.taskAdapter.adapt(item)))
    );
  }

  getTaskByNameAndVersion(name: string, version: string): Observable<Task> {
    return this.http.get(this.configService.getTaskByNameAndVersionURL(name, version)).pipe(
      map((item: any) => this.taskAdapter.adapt(item))
    );
  }

  searchTasksByName(name: string): Observable<Task[]> {
    return this.http.get(this.configService.searchTasksByNameURL(name)).pipe(
      map((data: any[]) => data.map((item: any) => this.taskAdapter.adapt(item)))
    );
  }

}
