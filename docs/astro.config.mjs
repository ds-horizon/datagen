import { defineConfig } from 'astro/config';
import starlight from '@astrojs/starlight';

// https://astro.build/config
export default defineConfig({
  site: 'https://ds-horizon.github.io/datagen/',
  base: '/datagen/',
  integrations: [
    starlight({
      title: 'datagen',
      description: 'Generate realistic data with a simple DSL',
      favicon: '/datagen-logo.png',
      logo: {
        src: './src/assets/datagen-logo.png',
      },
      social: [
        {
          icon: 'github',
          label: 'GitHub',
          href: 'https://github.com/ds-horizon/datagen',
        },
      ],
      sidebar: [
        {
          label: 'Introduction',
          items: [
            'introduction/overview',
            'introduction/getting-started',
          ],
        },
        {
          label: 'Concepts',
          items: [
            'concepts/data-model',
            {
              label: 'Advanced Concepts',
              items: [
                'concepts/advanced/overview',
                'concepts/advanced/optional-sections',
                'concepts/advanced/iter-variable',
                'concepts/advanced/function-calls',
                'concepts/advanced/model-references',
                'concepts/advanced/nesting',
              ],
            },
            {
              label: 'Sinks',
              items: [
                'sinks/overview',
                'sinks/config',
                'sinks/mysql',
              ],
            },
          ],
        },
        {
          label: 'CLI Reference',
          items: [
            'cli/datagenc-reference',
            'cli/datagen-reference',
          ],
        },
        {
          label: 'Examples',
          items: [
            {
              label: 'Fields',
              items: [
                'examples/1_fields/fields-overview',
                'examples/1_fields/single_field_model/single-field-model',
                'examples/1_fields/multi_field_model/multi-field-model',
              ],
            },
            'examples/2_calls/calls',
            'examples/3_misc/misc',
            'examples/4_iter/iter',
            'examples/5_reference/reference',
            {
              label: 'Metadata',
              items: [
                'examples/6_metadata/metadata-overview',
                'examples/6_metadata/count/count',
                'examples/6_metadata/tags/tags',
              ],
            },
          ],
        },
        {
          label: 'Language',
          items: [
            'language/dsl-specification',
            'language/built-in-functions',
          ],
        },
        {
          label: 'Reference',
          items: [
            'reference/datagenc-vs-datagen',
            'reference/troubleshooting',
          ],
        },
      ],
      customCss: ['./src/styles/custom.css'],
      tableOfContents: {
        minHeadingLevel: 2,
        maxHeadingLevel: 4,
      },
      pagination: true,
    }),
  ],
  server: {
    port: 4321,
    host: true,
  },
});
