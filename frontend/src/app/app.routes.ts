import type { Routes } from '@angular/router';

export const routes: Routes = [
  {
    path: '',
    // Lazy loading: each page ships in its own chunk, downloaded on demand.
    loadComponent: () => import('./features/home/pages/home/home.page').then((m) => m.HomePage),
  },
  {
    path: 'pantries/:slug',
    loadComponent: () =>
      import('./features/pantry/pages/pantry/pantry.page').then((m) => m.PantryPage),
  },
  { path: '**', redirectTo: '' },
];
