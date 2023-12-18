return {
	{
		"nvim-treesitter/nvim-treesitter",
		cmd = { "TSUpdateSync", "TSUpdate", "TSInstall" },
		build = ":TSUpdate",
		opts = {
			auto_install = true,
			highlight = {
				enable = true,
			},
			indent = {
				enable = true,
			},
			ensure_installed = {
				"vue",
				"typescript",
				"javascript",
				"tsx",
				"toml",
				"fish",
				"php",
				"json",
				"yaml",
				"css",
				"scss",
				"html",
				"lua",
				"python",
				"go",
				"markdown",
				"markdown_inline",
				"bash",
				"vim",
				"regex",
			},
			autotag = {
				enable = true,
			},
		},
		config = function(_, opts)
			require("nvim-treesitter.configs").setup(opts)
		end,
	},
}
