local icons = require("config.icons")

return {
	-- theme
	--"gruvbox-community/gruvbox",
  {
    "scottmckendry/cyberdream.nvim",
    lazy = false,
    priority = 1000,
    opts = {
      transparent = true,
      italic_comments = true,
    }
  },

	-- icons
	"kyazdani42/nvim-web-devicons",

	-- ui improvements
	{
		"stevearc/dressing.nvim",
		lazy = true,
		init = function()
			vim.ui.select = function(...)
				require("lazy").load({ plugins = { "dressing.nvim" } })
				return vim.ui.select(...)
			end
			vim.ui.input = function(...)
				require("lazy").load({ plugins = { "dressing.nvim" } })
				return vim.ui.input(...)
			end
		end,
	},

	-- notifications
	--[[{]]
		--[["rcarriga/nvim-notify",]]
		--[[opts = {]]
			--[[render = "wrapped-compact",]]
			--[[icons = {]]
				--[[DEBUG = icons.diagnostics.Debug,]]
				--[[ERROR = icons.diagnostics.Error,]]
				--[[INFO = icons.diagnostics.Info,]]
				--[[TRACE = icons.diagnostics.Trace,]]
				--[[WARN = icons.diagnostics.Warn,]]
			--[[},]]
			--[[timeout = 3000,]]
			--[[max_height = function()]]
				--[[return math.floor(vim.o.lines * 0.75)]]
			--[[end,]]
			--[[max_width = function()]]
				--[[return math.floor(vim.o.columns * 0.75)]]
			--[[end,]]
			--[[on_open = function(win)]]
				--[[vim.api.nvim_win_set_config(win, { zindex = 100 })]]
			--[[end,]]
		--[[},]]
		--[[init = function()]]
			--[[vim.notify = require("notify")]]
		--[[end,]]
	--[[},]]
}
