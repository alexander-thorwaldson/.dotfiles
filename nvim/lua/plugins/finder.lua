return {
	-- Aerial
	{
		"stevearc/aerial.nvim",
		opts = {},
	},

	-- Scope
	{
		"tiagovla/scope.nvim",
		opts = {},
	},

  -- Spectre
  {
    "nvim-pack/nvim-spectre",
    dependencies = {
      "nvim-lua/plenary.nvim",
    },
    cmd = "Spectre",
    opts = {
      is_open_target_win = true,
      is_insert_mode = true,
      line_sep = "-------------------------",
    },
  },

	-- Telescope
	{
		"nvim-telescope/telescope.nvim",
		dependencies = {
			"nvim-lua/plenary.nvim",
			"nvim-telescope/telescope-file-browser.nvim",
		},
		cmd = "Telescope",
		opts = {
			defaults = {
				sorting_strategy = "ascending",
				layout_config = {
					prompt_position = "top",
				},
			},
			extensions = {
				file_browser = {
					hijack_netrw = true,
					respect_gitignore = false,
					hidden = true,
					grouped = true,
					dir_icon = "",
					git_status = false,
					initial_mode = "normal",
				},
				repo = {
					list = {
						search_dirs = {
							"~/code",
						},
					},
				},
			},
		},
		config = function(_, opts)
			-- setup telescope
			local telescope = require("telescope")
			telescope.setup(opts)

			-- load extensions
			telescope.load_extension("aerial")
			telescope.load_extension("scope")
			telescope.load_extension("file_browser")
		end,
		keys = {
			{ "<Leader>ff", ":Telescope find_files<cr>" },
			{ "<Leader>fb", ":Telescope scope buffers<cr>" },
			{ "<Leader>ft", ":Telescope live_grep<cr>" },
			{ "<Leader>fa", ":Telescope aerial<cr>" },
			{ "<Leader>fd", ":Telescope file_browser<cr>" },
      { "<Leader>fe", ":Telescope diagnostics<cr>" }
		},
	},
}
