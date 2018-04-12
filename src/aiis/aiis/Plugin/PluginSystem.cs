using System;
using System.Collections;
using System.Collections.Generic;
using System.IO;
using System.Linq;
using System.Reflection;
using System.Text;
using System.Threading.Tasks;
using core;
using Nancy;

namespace aiis
{
    public class PluginInstance
    {
        public PluginInstance()
        {
            this.task = null;
        }

        public Type type;
        public Object instance;
        public Task task;
    }

    public class PluginModel
    {
        public string name;
        public string desc;
        public bool enable;

        public PluginInstance instance;
    }

    public class PluginSystem
    {
        public static List<PluginModel> GetPlugins()
        {
            lock(_plugins)
            {
                return _plugins;
            }
        }

        public static int Enable(int id, bool enable, ref PluginModel plugin)
        {
            int rt = 0;

            lock(_plugins)
            {
                PluginModel p = null;
                try
                {
                    p = _plugins[id - 1];
                }
                catch(Exception e)
                {
                    // 找不到
                    rt = 1;
                    return rt;
                }

                if(p.enable == enable)
                {
                    return 0;
                }
                else
                {
                    p.enable = enable;
                    if(p.enable == true)
                    {
                        // 启用插件
                        p.instance.task = Task.Run(() =>
                        {
                            var type = p.instance.type;
                            var instance = p.instance.instance;
                            type.GetMethod(PluginConfig.FUNC_RUN).Invoke(instance, null);
                        });
                    }
                    else
                    {
                        // 停用插件
                        p.instance.type.GetMethod(PluginConfig.FUNC_RELEASE).Invoke(p.instance.instance, null);
                    }

                    _plugins[id - 1] = p;
                }

                plugin = p;
            }

            return 0;
        }

        public static void LoadFromConfig(AiisConfig config)
        {
            for(int i = 0; i < config.plugins.Count; ++i)
            {
                try
                {
                    PluginModel plugin = new PluginModel();
                    plugin.instance = new PluginInstance();

                    // 装载插件
                    string plugin_path = Path.Combine(plugin_sub_path, config.plugins[i].path);
                    Assembly assembly = Assembly.LoadFrom(plugin_path);

                    plugin.instance.type = assembly.GetType(PluginConfig.DEFAULT_NAME);
                    plugin.instance.instance = Activator.CreateInstance(plugin.instance.type, null);

                    plugin.name = plugin.instance.type.GetMethod(PluginConfig.FUNC_NAME).Invoke(plugin.instance.instance, null).ToString();
                    plugin.desc = plugin.instance.type.GetMethod(PluginConfig.FUNC_DESC).Invoke(plugin.instance.instance, null).ToString();
                    plugin.enable = config.plugins[i].enable;

                    if (plugin.enable)
                    {
                        // 运行插件
                        plugin.instance.task = Task.Run(() =>
                        {
                            var type = plugin.instance.type;
                            var instance = plugin.instance.instance;
                            type.GetMethod(PluginConfig.FUNC_RUN).Invoke(instance, null);
                        });
                    }

                    _plugins.Add(plugin);

                    Console.WriteLine(String.Format("init plugin [{0}] OK", config.plugins[i].path));
                }
                catch (Exception e)
                {
                    Console.WriteLine(String.Format("init plugin [{0}] FAILED:{1}", config.plugins[i].path, e.Message));
                    continue;
                }
            }
            
        }
        
        private static List<PluginModel> _plugins = new List<PluginModel>();
        private const string plugin_sub_path = "./plugins";
    }
}
