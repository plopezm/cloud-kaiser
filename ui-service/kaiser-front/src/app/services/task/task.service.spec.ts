import {TestBed} from '@angular/core/testing';

import {TaskService} from './task.service';
import {HttpClientModule} from "@angular/common/http";

describe('TaskService', () => {
  beforeEach(() => TestBed.configureTestingModule({
    imports: [
      HttpClientModule
    ]
  }));

  it('should be created', () => {
    const service: TaskService = TestBed.get(TaskService);
    expect(service).toBeTruthy();
  });
});
