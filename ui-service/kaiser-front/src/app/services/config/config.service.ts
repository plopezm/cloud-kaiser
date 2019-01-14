import {Injectable} from "@angular/core";
import {environment} from "../../../environments/environment";

@Injectable({
  providedIn: "root"
})
export class ConfigService {

  constructor() {}

  getApiEndpoint() {
    return environment.apiEndpoint;
  }

  getWSEndpoint() {
    return environment.wsEndpoint;
  }

  getTasksURL(): string {
    return `${environment.apiEndpoint}/query/v1/tasks`;
  }

  getTasksByNameURL(name: string): string {
    return `${environment.apiEndpoint}/query/v1/tasks/${name}`;
  }

  getTaskByNameAndVersionURL(name: string, version: string): string {
    return `${environment.apiEndpoint}/query/v1/tasks/${name}/${version}`;
  }

  searchTasksByNameURL(name: string): string {
    return `${environment.apiEndpoint}/query/v1/search/tasks?query=${name}`;
  }

}
