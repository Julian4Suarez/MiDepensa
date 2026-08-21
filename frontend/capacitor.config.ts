import type { CapacitorConfig } from '@capacitor/cli';

// Native platforms are not generated yet. Run `npx cap add android` when you
// want to build the mobile app; this config is already wired for it.
const config: CapacitorConfig = {
  appId: 'work.midepensa.app',
  appName: 'MiDepensa',
  webDir: 'www',
};

export default config;
