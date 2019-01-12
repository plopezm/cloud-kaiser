import {Task} from "./task.model";
import {Injectable} from "@angular/core";
import {Adapter} from "../../adapter";

@Injectable({
  providedIn: "root"
})
export class TaskAdapter implements Adapter<Task> {
  adapt(item: any): Task {
    return new Task(item.name, item.version, item.script, item.createdAt);
  }
}
