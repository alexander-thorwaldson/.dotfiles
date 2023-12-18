return {
	{
		"neovim/nvim-lspconfig",
		config = function(_, _)
			local lspconfig = require("lspconfig")
			local capabilities = require("cmp_nvim_lsp").default_capabilities()
			lspconfig.tsserver.setup({
				capabilities = capabilities,
        init_options = {
          plugins = {
            {
              name = "@vue/typescript-plugin",
              location = "/usr/lib/node_modules/@vue/typescript-plugin",
              languages = {"javascript", "typescript", "vue"},
            },
          },
        },
        filetypes = {
          "javascript",
          "typescript",
          "vue",
        },
			})
      lspconfig.volar.setup({
				capabilities = capabilities,
			})
			lspconfig.pyright.setup({
				capabilities = capabilities,
			})
			lspconfig.gopls.setup({
				capabilities = capabilities,
			})
			lspconfig.lua_ls.setup({
				capabilities = capabilities,
			})
		end,
	},
	{
		"https://git.sr.ht/~whynothugo/lsp_lines.nvim",
		opts = {},
	},
}
