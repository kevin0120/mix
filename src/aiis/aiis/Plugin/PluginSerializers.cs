using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace aiis
{
    public class PluginSerializer
    {
        public int id;
        public string name;
        public string desc;
        public bool enable;
        public string status;
    }

    public class PluginSerializers
    {
        public static PluginSerializer PluginItemSerializer(PluginModel plugin)
        {
            var ps = new PluginSerializer();

            ps.name = plugin.name;
            ps.desc = plugin.desc;
            ps.enable = plugin.enable;
            ps.status = "Stopped";

            if (ps.enable)
            {
                if (plugin.instance.task != null)
                {
                    if (!plugin.instance.task.IsCompleted)
                    {
                        ps.status = "Running";
                    }
                }
            }

            return ps;
        }

        public static List<PluginSerializer> PluginListSerializer(List<PluginModel> plugins)
        {
            var rt = new List<PluginSerializer>();

            for (int i = 0; i < plugins.Count; ++i)
            {
                PluginSerializer ps = new PluginSerializer();
                ps.id = i + 1;
                ps.name = plugins[i].name;
                ps.desc = plugins[i].desc;
                ps.enable = plugins[i].enable;
                ps.status = "Stopped";

                if(ps.enable)
                {
                    if(plugins[i].instance.task != null)
                    {
                        if(!plugins[i].instance.task.IsCompleted)
                        {
                            ps.status = "Running";
                        }
                    }
                }

                rt.Add(ps);
            }

            return rt;
        }
    }
}
