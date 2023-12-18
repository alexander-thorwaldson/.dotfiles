local colors = {
	bg = "None",
	bg1 = "#1e2124",
	bg2 = "#3c4048",
	gray = "#7b8496",
	fg = "#ffffff",
	yellow = "#f1ff5e",
	cyan = "#5ef1ff",
	blue = "#5ea1ff",
	green = "#5eff6c",
	orange = "#ffbd5e",
	violet = "#bd5eff",
	magenta = "#ff5ef1",
	pink = "#ff5ea0",
	red = "#ff6e5e",
}

-- Selects a color based on neovims mode
local function mode_color()
	local modes = {
		n = colors.violet,
		i = colors.orange,
		v = colors.blue,
		[""] = colors.blue,
		V = colors.blue,
		c = colors.cyan,
		no = colors.magenta,
		s = colors.blue,
		S = colors.blue,
		[""] = colors.blue,
		ic = colors.yellow,
		R = colors.violet,
		Rv = colors.violet,
		cv = colors.red,
		ce = colors.red,
		r = colors.cyan,
		rm = colors.cyan,
		["r?"] = colors.cyan,
		["!"] = colors.red,
		t = colors.red,
	}
	return modes[vim.fn.mode()]
end

return {
	-- Commenter
	"preservim/nerdcommenter",

	-- Lualine
	{
		"nvim-lualine/lualine.nvim",
		opts = {
			options = {
				-- Disable sections and component separators
				component_separators = "",
				section_separators = "",
				theme = {
					normal = { c = { fg = colors.fg, bg = colors.bg1 } },
					inactive = { c = { fg = colors.fg, bg = colors.bg1 } },
				},
			},
			sections = {
				-- inactive sections
				lualine_a = {},
				lualine_b = {},
				lualine_y = {},
				lualine_z = {},
				-- active sections
				lualine_c = {
					{
						"mode",
						color = function()
							return { bg = mode_color(), gui = "bold" }
						end,
						padding = { left = 2, right = 2 },
					},
					{
						"filename",
						color = { fg = colors.magenta, bg = colors.bg2, gui = "bold" },
					},
					{
						"diagnostics",
						sources = { "nvim_diagnostic" },
						symbols = { error = " ", warn = " ", info = " " },
						diagnostics_color = {
							color_error = { fg = colors.red },
							color_warn = { fg = colors.yellow },
							color_info = { fg = colors.cyan },
						},
						color = { bg = colors.bg2 },
					},
					{ "progress" },
					{ "location" },
				},
				lualine_x = {
					{
						"o:encoding",
						fmt = string.upper,
						color = { fg = colors.fg },
					},
					{
						"fileformat",
						fmt = string.upper,
						icons_enabled = false,
						color = { fg = colors.fg },
					},
					{ "filetype" },
					{
						"diff",
						symbols = { added = " ", modified = " ", removed = " " },
						diff_color = {
							added = { fg = colors.green },
							modified = { fg = colors.orange },
							removed = { fg = colors.red },
						},
						color = { bg = colors.bg2 },
					},
					{ "branch", color = { fg = colors.magenta, bg = colors.bg2, gui = "bold" } },
					{
						"datetime",
						style = "%H:%M:%S",
						color = function()
							return { bg = mode_color(), gui = "bold" }
						end,
						padding = { left = 2, right = 2 },
					},
				},
			},
			tabline = {
				lualine_c = {
					{
						function()
							return ""
						end,
						color = function()
							return { bg = mode_color(), gui = "bold" }
						end,
						padding = { left = 2, right = 2 },
					},
					{
						"buffers",
						icons_enabled = false,
						mode = 2,
						buffers_color = {
							active = { fg = colors.magenta, gui = "bold" },
							inactive = { bg = colors.bg2, fg = colors.gray },
						},
						symbols = {
							modified = " +",
							alternate_file = "",
						},
						color = { bg = colors.bg2 },
					},
				},
				lualine_x = {
					{
						"tabs",
						show_modified_status = false,
						tabs_color = {
							active = { fg = colors.magenta, bg = colors.bg2, gui = "bold" },
							inactive = { fg = colors.gray, bg = colors.bg2 },
						},
					},
					{
						function()
							return "project"
						end,
						color = function()
							return { bg = mode_color(), gui = "bold" }
						end,
						padding = { left = 2, right = 2 },
					},
				},
			},
			inactive_sections = {
				-- remove defaults
				lualine_a = {},
				lualine_b = {},
				lualine_y = {},
				lualine_z = {},
				lualine_c = {},
				lualine_x = {},
			},
		},
	},

	-- indent guides for Neovim
	{
		"lukas-reineke/indent-blankline.nvim",
		opts = {
			indent = {
				char = "│",
				tab_char = "│",
			},
			scope = { enabled = false },
			exclude = {
				filetypes = {
					"help",
					"alpha",
					"dashboard",
					"neo-tree",
					"Trouble",
					"trouble",
					"lazy",
					"mason",
					"notify",
					"toggleterm",
					"lazyterm",
				},
			},
		},
		main = "ibl",
	},

	{
		"echasnovski/mini.indentscope",
		version = false,
		opts = function()
			return {
				symbol = "│",
				options = { try_as_border = true },
				draw = {
					delay = 0,
					animation = require("mini.indentscope").gen_animation.none(),
				},
			}
		end,
		init = function()
			vim.api.nvim_create_autocmd("FileType", {
				pattern = {
					"help",
					"alpha",
					"dashboard",
					"neo-tree",
					"Trouble",
					"trouble",
					"lazy",
					"mason",
					"notify",
					"toggleterm",
					"lazyterm",
				},
				callback = function()
					vim.b.miniindentscope_disable = true
				end,
			})
		end,
	},
}
