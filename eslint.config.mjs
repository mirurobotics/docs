import * as mdx from "eslint-plugin-mdx";

export default [
  {
    ignores: [
      ".ai/**",
      "api/**",
      "images/**",
      "logo/**",
      "node_modules/**",
      "output/**",
      "videos/**",
      "changelog/**"
    ]
  },
  {
    ...mdx.flat,
    files: ["docs/**/*.mdx", "tests/lint-fixtures/**/*.mdx"],
    rules: {
      ...mdx.flat.rules,
      "mdx/remark": "off"
    }
  }
];
