import { ChangeDetectionStrategy, Component, computed, inject, signal } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { Router } from '@angular/router';
import {
  IonButton,
  IonContent,
  IonIcon,
  IonInput,
  IonNote,
  IonSpinner,
} from '@ionic/angular/standalone';
import { addIcons } from 'ionicons';
import { basketOutline } from 'ionicons/icons';
import { firstValueFrom } from 'rxjs';

import { PantryApiService } from '../../../../core/services/pantry-api.service';
import { slugify } from '../../../../core/utils/slugify';

/** Landing page: describes the app and creates a pantry from a name. */
@Component({
  selector: 'app-home',
  standalone: true,
  changeDetection: ChangeDetectionStrategy.OnPush,
  imports: [FormsModule, IonButton, IonContent, IonIcon, IonInput, IonNote, IonSpinner],
  templateUrl: './home.page.html',
  styleUrl: './home.page.scss',
})
export class HomePage {
  private readonly api = inject(PantryApiService);
  private readonly router = inject(Router);

  protected readonly name = signal('');
  protected readonly submitting = signal(false);
  protected readonly error = signal<string | null>(null);

  /** Live preview of the URL the pantry will get. */
  protected readonly slug = computed(() => slugify(this.name()));
  protected readonly canSubmit = computed(() => this.slug().length > 0 && !this.submitting());

  constructor() {
    addIcons({ basketOutline });
  }

  protected async create(): Promise<void> {
    if (!this.canSubmit()) {
      return;
    }

    this.submitting.set(true);
    this.error.set(null);
    try {
      const pantry = await firstValueFrom(this.api.createPantry(this.name()));
      await this.router.navigate(['/pantries', pantry.slug]);
    } catch (error: unknown) {
      this.error.set(
        (error as { status?: number }).status === 409
          ? 'That name is taken. Try adding a surname or a number.'
          : 'The pantry could not be created. Is the API running?',
      );
    } finally {
      this.submitting.set(false);
    }
  }
}
