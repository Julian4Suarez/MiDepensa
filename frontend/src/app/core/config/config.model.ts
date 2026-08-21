/** Shape of the runtime configuration served at /config.json. */
export interface AppConfig {
  /** Absolute base URL of the API, including the version prefix. */
  BACKEND_URL: string;
}
