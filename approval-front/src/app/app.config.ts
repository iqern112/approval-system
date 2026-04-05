// src/app/app.config.ts
import { ApplicationConfig } from '@angular/core';
import { provideHttpClient } from '@angular/common/http'; // เพิ่มบรรทัดนี้
import { provideRouter } from '@angular/router';
import { routes } from './app.routes';

export const appConfig: ApplicationConfig = {
  providers: [
    provideHttpClient(), // ใส่เข้าไปในลิสต์นี้
    provideRouter(routes)
  ]
};