import {Component, OnInit, ViewChild} from '@angular/core';
import {TaskService} from "../../services/task/task.service";
import {Task} from "../../core/models/task/task.model";
import {MatPaginator, MatTableDataSource} from "@angular/material";
import {WebSocketService} from "../../services/wsocket/wsocket.service";
import {TOPIC_TASK_CREATED} from "../../services/wsocket/topics";

@Component({
  selector: 'app-management',
  templateUrl: './management.component.html',
  styleUrls: ['./management.component.css']
})
export class ManagementComponent implements OnInit {
  displayedColumns: string[] = ['version', 'name', 'script', 'createdAt'];
  @ViewChild(MatPaginator) paginator: MatPaginator;

  tasksDatasource: MatTableDataSource<Task>;

  constructor(private taskService: TaskService, private webSocketService: WebSocketService) {
  }

  ngOnInit() {
    this.tasksDatasource = new MatTableDataSource();
    this.taskService.getAllTasks(0, 20).subscribe(value => this.tasksDatasource.data = value);
    this.webSocketService.onEvent(TOPIC_TASK_CREATED).subscribe(msg => {
      let rows = [...this.tasksDatasource.data, msg.content];
      this.tasksDatasource.data = rows;
    });
  }

}
