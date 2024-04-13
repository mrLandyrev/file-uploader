import { bootstrapApplication } from '@angular/platform-browser';
import { appConfig } from './modules/app/app.config';
import { AppComponent } from './modules/app/app.component';

bootstrapApplication(AppComponent, appConfig)
  .catch((err) => console.error(err));
