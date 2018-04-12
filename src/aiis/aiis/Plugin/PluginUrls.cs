using Nancy;
using Newtonsoft.Json;
using System;
using System.Collections.Generic;
using System.Linq;
using System.Text;
using System.Threading.Tasks;

namespace aiis.Plugin
{
    public class EnablePlugin
    {
        public bool enable;
    }

    public class PluginUrls : NancyModule
    {
        public PluginUrls()
        {
            // 插件列表
            Get[ApiConfig.API_PREFIX() + "/plugins"] = _ =>
            {
                return Response.AsJson(PluginSerializers.PluginListSerializer(PluginSystem.GetPlugins()));
            };

            // 插件启用停用
            Patch[ApiConfig.API_PREFIX() + "/plugins/{id:int}"] = parameters =>
            {
                EnablePlugin ep = null;
                try
                {
                    ep = JsonConvert.DeserializeObject<EnablePlugin>(this.Request.Body.ToString());
                }
                catch(Exception e)
                {
                    return ResponseExtensions.WithStatusCode(Response.AsText("request error"), HttpStatusCode.BadRequest);
                }

                PluginModel pm = null;
                int rt = PluginSystem.Enable(parameters.id, ep.enable, ref pm);
                if(rt == 1)
                {
                    return ResponseExtensions.WithStatusCode(Response.AsText("plugin not found"), HttpStatusCode.NotFound);
                }

                PluginSerializer ps = PluginSerializers.PluginItemSerializer(pm);
                ps.id = parameters.id;
                return Response.AsJson(ps);
            };
        }
    }
}
