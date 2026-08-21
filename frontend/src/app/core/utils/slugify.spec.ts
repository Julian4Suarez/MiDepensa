import { slugify } from './slugify';

describe('slugify', () => {
  it.each([
    ['Familia Suarez', 'familia-suarez'],
    ['Familia Suárez', 'familia-suarez'],
    ['La Peña', 'la-pena'],
    ['casa   ---   playa', 'casa-playa'],
    ['casa/../etc?x=1', 'casa-etc-x-1'],
    ['  --Casa--  ', 'casa'],
    ['Piso 2', 'piso-2'],
  ])('turns %s into %s', (input, expected) => {
    expect(slugify(input)).toBe(expected);
  });

  it('returns an empty string when there is nothing URL-safe', () => {
    expect(slugify('///')).toBe('');
    expect(slugify('   ')).toBe('');
  });

  it('truncates to 60 characters without a trailing hyphen', () => {
    const result = slugify('aaaaaaaaaa bbbbbbbbbb cccccccccc dddddddddd eeeeeeeeee ffffffffff gggg');

    expect(result.length).toBeLessThanOrEqual(60);
    expect(result.endsWith('-')).toBe(false);
  });
});
