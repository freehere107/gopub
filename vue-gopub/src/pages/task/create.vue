<template>
  <div>

    <div class="panel">
      <panel-title :title="$route.meta.title"></panel-title>
      <div class="panel-body"
           v-loading="load_data"
           element-loading-text="拼命加载中">
        <el-row>
          <el-col :span="20">
            <el-form label-width="100px">
              <div class="panel-title">
                预发布
                <div class="fr">

                </div>
              </div>
              <div class="panel-body">
                <el-form-item label="项目名称:" label-width="100px">
                  <el-select v-model.number="pro_id" filterable @change="select_project" placeholder="请选择"
                             style="width: 400px;">
                    <el-option
                      v-for="item in options"
                      :key="item.value"
                      :label="item.label"
                      :value="item.value">
                      <span style="float: left">{{ item.label }}</span>
                      <span style="float: right; color: #8492a6; font-size: 13px">{{ item.lockstatus }}</span>
                    </el-option>
                  </el-select>

                </el-form-item>
              </div>
              <el-form-item>
                <el-button type="primary" @click="on_submit_form" :loading="on_submit_loading"
                           :disabled="btn_submit_disable">创建
                </el-button>
                <el-button @click="add_lock" :loading="on_submit_loading" :disabled="btn_add_lock_disable">锁定
                </el-button>
                <el-button @click="free_lock" :loading="on_submit_loading" :disabled="btn_free_lock_disable">解锁
                </el-button>

                <el-button @click="$router.back()">返回</el-button>
              </el-form-item>
            </el-form>
          </el-col>
        </el-row>
      </div>
    </div>
  </div>
</template>
<script type="text/javascript">
  import {panelTitle} from 'components'
  import {port_conf} from 'common/port_uri'
  import store from 'store'


  export default {
    data() {
      return {
        projects: null,
        options: [],
        pro_id: null,
        load_data: false,
        on_submit_loading: false,
        btn_submit_disable: true,
        btn_add_lock_disable: true,
        btn_free_lock_disable: true
      }
    },
    created() {
      this.get_project_data()
    },
    methods: {
      //获取数据
      get_project_data() {
        this.load_data = true
        this.$http.post(port_conf.mylist, {page: 1, length: 100})
          .then(({data: {data}}) => {
            let opData = []
            let uid = store.state.user_info.user.Id
            for (let i in data.table_data) {
              let value = data.table_data[i].id
              let env = ""
              let locks_tatus = "";
              if (data.table_data[i].user_lock > 0) {
                if (data.table_data[i].user_lock == uid) {
                  locks_tatus = data.table_data[i].lockuser + "锁定中"
                } else {
                  locks_tatus = "锁定中"
                }
              }

              env = "测试环境"
              let lable = env + "-" + data.table_data[i].name
              opData.push({label: lable, value: value, lockstatus: locks_tatus})
            }
            this.projects = data.table_data
            this.options = opData
            this.load_data = false
            this.select_project()
            if (this.$route.query.id) {
              this.pro_id = this.$route.query.id
              this.on_submit_form()
            }
          }).catch(() => {
          this.load_data = false
        })
      },
      add_lock() {
        let proId = 0

        proId = this.pro_id

        if (proId) {
          this.$http.post(port_conf.lock, {
            project_id: proId,
            act: 1
          }).then(({data: {data}}) => {
            this.$message({
              message: "锁定成功！",
              type: 'success'
            })
            this.get_project_data()
          })
        }
      },
      free_lock() {
        let proId = 0

        proId = this.pro_id

        if (proId) {
          this.$http.post(port_conf.lock, {
            project_id: proId,
            act: 0
          }).then(({data: {data}}) => {
            this.$message({
              message: "解锁成功！",
              type: 'success'
            })
            this.get_project_data()
          })
        }
      },
      select_project() {
        let uid = store.state.user_info.user.Id
        let role = store.state.user_info.user.Role
        let proId = 0

        proId = this.pro_id
        if (!proId) {
          return
        }

        for (let i in this.projects) {
          let project = this.projects[i]
          if (project.id == proId) {
            if (project.user_lock > 0) {
              if (project.user_lock == uid) {
                this.btn_submit_disable = false
                this.btn_add_lock_disable = true
                this.btn_free_lock_disable = false
              } else {
                this.btn_submit_disable = true
                this.btn_add_lock_disable = true
                if (role == 1) {
                  this.btn_free_lock_disable = false
                } else {
                  this.btn_free_lock_disable = true
                }
              }
            } else {
              this.btn_submit_disable = false
              this.btn_add_lock_disable = false
              this.btn_free_lock_disable = true
            }
          }
        }
      },
      //提交
      on_submit_form() {
        let proId = 0
        proId = this.pro_id

        if (proId) {
          for (let i in  this.projects) {
            let project = this.projects[i]
            if (proId == project.id) {
              if (project.repo_type === "git") {
                this.$router.push({
                  name: 'taskGit',
                  query: {id: proId}
                })
                return
              }

            }
          }

        }
      }
    },
    components: {
      panelTitle
    }
  }
</script>
